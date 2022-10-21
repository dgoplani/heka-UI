import { Component, ComponentRef, Input, OnInit, ViewChild, ViewContainerRef } from '@angular/core';
import { animate, state, style, transition, trigger } from '@angular/animations';
import { UiNotificationComponent } from '../ui-notification/ui-notification.component';
import { Notification } from 'src/app/interfaces/notification';
import { NotificationService } from 'src/app/ui-services/notification.service';
import { UiModalComponent } from '../ui-modal/ui-modal.component';
import { NodeData } from 'src/app/interfaces/node-data';

@Component({
  selector: 'app-ui-notify-plane',
  templateUrl: './ui-notify-plane.component.html',
  styleUrls: ['./ui-notify-plane.component.css'],
  animations: [
    trigger('slide', [
      state('true', style({width: '270px', opacity: 1, display: 'inline'})),
      state('false', style({width: '0px', opacity: 0, display: 'none'})),
      transition('0 <=> 1', [
        style({ display: 'inline' }), 
        animate('200ms ease')
      ])
    ])
  ]
})

export class UiNotifyPlaneComponent implements OnInit {
  
  @Input('node_list')nodeList!: NodeData[];

  // Hotfix Modal To display Hotfix Infomation
  @ViewChild('hotfixModal', {static: false})
  hotfixModal!: UiModalComponent;
  // Active Notification Component Reference
  @ViewChild('notification', { read: ViewContainerRef })
  nAlertPlane!: ViewContainerRef;
  // History Notification Component Reference
  @ViewChild('notification_history', { read: ViewContainerRef })
  nHistoryAlertPlane!: ViewContainerRef;
  
  lastHotfixNotficationTag!: string;
  isOpen: boolean = false;
  lastNotificationTS: number = 0;
  stdDelay: number = 200; // 0.2sec
  stdExitDelay: number = 10000;  // 10sec
  sequenceNotificationCounter: number = 0;
  constructor(private notificationService: NotificationService) { }

  ngOnInit(): void { 
    this.notificationService.notifications.subscribe((ndata: Notification) => {
      //console.log(ndata);
      if(ndata.source == 'HOTFIX' && this.lastHotfixNotficationTag != ndata.tag) {
        console.log('Hotfix Tag Changed from ', this.lastHotfixNotficationTag, ' to ', ndata.tag);
        this.lastHotfixNotficationTag = ndata.tag;
        this.nAlertPlane.clear();
        this.nHistoryAlertPlane.clear();
      }
      this.handleNotification(ndata);
    });
  }

  handleNotification(ndata: Notification) {
    switch(ndata.source){
      case 'HOTFIX':
        this.handleHotfixNotification(ndata);
        break;
      case 'CHILD_NOTIFICATION':
        this.handleChildNotification(ndata);
        break;
      default:
        break;
    }
  }

  calculateDelay(): number {
        let timeout = 0;
        let ts = Date.now();
        let diff_ts = ts - this.lastNotificationTS;
        this.lastNotificationTS = ts;
        //console.log("diff: ", diff_ts);
        if(diff_ts < this.stdDelay){
          timeout = this.stdDelay - diff_ts;
          this.sequenceNotificationCounter++;
        } else {
          this.sequenceNotificationCounter = 0;
        }
        //console.log("timeout: ", timeout);
        let delay = (this.sequenceNotificationCounter-1) * this.stdDelay + timeout;
        //console.log("delay: ", delay);
        return delay;
  }

  handleHotfixNotification(ndata: Notification) {
    let delay = this.calculateDelay();
    setTimeout(() => {
      let n = this.nAlertPlane.createComponent<UiNotificationComponent>(UiNotificationComponent);
      n.instance.initNotification(ndata, 'ALERT_PLANE', n);
      n.instance.showNotification(this.stdExitDelay + delay);
    }, delay);
    this.addToHistoryNotification(ndata);
  }

  addToHistoryNotification(ndata: Notification) {
    let nh = this.nHistoryAlertPlane.createComponent<UiNotificationComponent>(UiNotificationComponent);
    nh.instance.initNotification(ndata, 'HISTORY_PLANE', nh);
    nh.instance.showNotification(0);
  }

  handleChildNotification(ndata: Notification) {
    switch(ndata.type){
      case 'TASK':
        switch(ndata.action){
          case 'SHOW':
            this.hotfixModal.openModal(ndata.extraData, this.nodeList);
            break;
          default:
            break;
        }
        break;
      default:
        break;
    }
  }

  open() {
    this.isOpen = true;
    this.nAlertPlane.clear();
  }

  close() {
    this.isOpen = false;
  }

  clearAll() {
    this.nHistoryAlertPlane.clear();
  }
}
