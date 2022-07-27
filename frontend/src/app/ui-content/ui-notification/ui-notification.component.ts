import { AfterViewInit, Component, ComponentRef, Directive, ElementRef, OnDestroy, OnInit, Renderer2, ViewChild } from '@angular/core';
import { Notification } from 'src/app/interfaces/notification';
import { NotificationService } from 'src/app/ui-services/notification.service';


@Component({
  selector: 'app-ui-notification',
  templateUrl: './ui-notification.component.html',
  styleUrls: ['./ui-notification.component.css']
})
export class UiNotificationComponent implements OnInit, OnDestroy, AfterViewInit {

  @ViewChild('notificationElement')nEleRef!: ElementRef;
  _ref!: ComponentRef<UiNotificationComponent>;
  ndata!: Notification;

  parent!: string;
  timer!: any;
  timeout: number = 0;
  visible: boolean = false;
  constructor(private notification_service: NotificationService, private renderer: Renderer2) { }
  ngAfterViewInit(): void {
    // this.renderer.listen(this.nEleRef.nativeElement, 'animationstart', (e) => {
    //   console.log('Animation Start');
    // });
    this.renderer.listen(this.nEleRef.nativeElement, 'animationend', (e) => {
      if (!this.visible) {
        this._ref.destroy();
      }
    });
  }

  ngOnDestroy(): void {
    console.log("Notification: I am Destroyed");
  }

  ngOnInit(): void {
  }

  initNotification(data: Notification, parent: string, selfRef: ComponentRef<UiNotificationComponent>) {
    this.ndata = data;
    this.parent = parent;
    this._ref = selfRef;
  }

  stopTimer() {
    clearTimeout(this.timer);
  }

  restartTimer() {
    if(this.timeout > 0) {
      this.timer = setTimeout(()=> {
        this.hideNotification();
      }, this.timeout/2);
    }
  }

  showNotification(exitTimeout: number): void {
    this.visible = true;
    if(exitTimeout > 0) {
      this.timeout = exitTimeout;
      this.timer = setTimeout(()=> {
        this.hideNotification();
      }, exitTimeout);
    }
  }

  hideNotification() {
    this.visible = false;
  }

  notificationAction() {
    switch(this.ndata.source) {
      case 'HOTFIX':
        this.hotfixAction();
        break;
      default:
        break;
    }
  }

  hotfixAction() {
    let nd: Notification = {
      source: 'CHILD_NOTIFICATION',
      type: 'TASK',
      tag: this.parent,
      title: '',
      body: '',
      action: 'SHOW',
      extraData: this.ndata.extraData,
    }
    this.notification_service.sendNotification(nd);
  }

}

