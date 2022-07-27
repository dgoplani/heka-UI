import { Component, Input, OnInit } from '@angular/core';
import { Hotfix, HotfixRequiredActions } from 'src/app/interfaces/cloud-hotfix-manifest';
import { NodeData } from 'src/app/interfaces/node-data';

@Component({
  selector: 'app-ui-modal',
  templateUrl: './ui-modal.component.html',
  styleUrls: ['./ui-modal.component.css']
})
export class UiModalComponent implements OnInit {

  @Input('plane_level') z_index!: number;
  isOpen = false;
  hf!: Hotfix;
  nodes!: NodeData[];

  constructor() { }

  ngOnInit(): void {
  }

  openModal(hf_data: Hotfix, node_list:NodeData[]) {
    //console.log(hf_data, node_list);
    this.hf = hf_data;
    this.nodes = node_list;
    if(this.nodes != undefined && this.hf.compatibleNode != 'ALL'){
      this.getCompatibleNodes();
    }
    this.isOpen = true;
  }

  closeModal() {
    this.isOpen = false;
  }

  getRequiredActionString(action: HotfixRequiredActions): string {
    if(action.systemReboot.toLowerCase().match('yes')){
      return 'System Reboot';
    } else if (action.productRestart.toLowerCase().match('yes')) {
      return 'Product Restart';
    } else if (action.serviceRestart.length > 0) {
      return 'Service Restart(' + action['serviceRestart'].join(', ') + ')';
    }
    return 'Not Available';
  }

  getCompatibleNodes() {
    this.nodes = this.nodes.filter((n:any) => n.role === this.hf.compatibleNode);
  }

  openLink(link: any) {
    window.open(link);
  }
}
