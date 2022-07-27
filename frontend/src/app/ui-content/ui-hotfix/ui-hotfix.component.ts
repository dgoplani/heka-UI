import { Component, ElementRef, Input, OnInit, ViewChild } from '@angular/core';
import { HotfixService } from 'src/app/ui-services/hotfix.service';
import { Notification } from 'src/app/interfaces/notification';
import { NotificationService } from 'src/app/ui-services/notification.service';
import { UiService } from 'src/app/ui-services/ui.service';
import { NodeData } from 'src/app/interfaces/node-data';
import { NodeHotfixData } from 'src/app/interfaces/node-hotfix-data';
import { CloudHotfixManifest, Hotfix, HotfixFixes, HotfixRequiredActions, ProcessedHotfix } from 'src/app/interfaces/cloud-hotfix-manifest';

@Component({
  selector: 'app-ui-hotfix',
  templateUrl: './ui-hotfix.component.html',
  styleUrls: ['./ui-hotfix.component.css']
})
export class UiHotfixComponent implements OnInit {

  @ViewChild('search')search!: ElementRef;
  @ViewChild('searchNode')searchNode!: ElementRef;
  @ViewChild('nodeListMenu')nodeListMenu!: ElementRef;
  
  hotfix_info!: CloudHotfixManifest;
  node_loading: boolean = false;
  hotfix_loading: boolean = false;
  loading_title: string = "Please wait a moment...";
  allNodes: NodeData[] = [];
  listNodes: NodeData[] = [];
  selectedNode!: NodeData;
  selectedNodeData!: NodeHotfixData;
  
  data: ProcessedHotfix[] = [];
  displayData: ProcessedHotfix[] = [];

  displayCols: string[] = ['applyStatus', 'name', 'severity', 'type', 'released', 'impactedArea', 'requiredActions'];
  displayColsMap: {[key:string]:string} = {
    'applyStatus': 'Status', 
    'name': 'Hotfix Name',
    'severity': 'Severity',
    'type': 'Type',
    'released': 'Release Date',
    'impactedArea': 'Impacted Area(s)',
    'requiredActions': 'Action Required'
  };
  
  hotfixSummary: {[key:string]:{[key:string]:number}} = {
    'IMPORTANT' : { 'available': 0, 'installed': 0 },
    'RECOMMENDED' : { 'available': 0, 'installed': 0 },
    'OPTIONAL' : { 'available': 0, 'installed': 0 }
  }

  sortColumn: string = '';
  sortOrder: string = '';
  filterData:Map<string, Map<string, boolean>> = new Map<string, Map<string, boolean>>();
  filterApplied: boolean = false;
  searchApplied: boolean = false;
  foundCount: number = 0;
  filterColumns: string[] = ['applyStatus', 'severity', 'type', 'impactedArea', 'requiredActions'];

  constructor(private ui_service: UiService,
              private hf_service: HotfixService, 
              private notify_service: NotificationService) { }

  ngOnInit(): void { 
    this.node_loading = true;
    this.hotfix_loading = true;
    this.ui_service.getNodeList().subscribe({
      next: (rdata: NodeData[]) => {
        this.allNodes = rdata;
        this.listNodes = this.allNodes;
        this.selectedNode = this.listNodes.filter((ele: NodeData) => ele.role == "MASTER")[0];
        console.log(this.listNodes);
        console.log(this.selectedNode);
      },
      error: (err: any) => {
        this.hf_service.handleError(err);
      },
      complete: () => {
        console.log('Got Node_List');
        this.hf_service.getNodeInfo(this.selectedNode.unique_id).subscribe({
          next: (rdata: NodeHotfixData) => {
            this.selectedNodeData = rdata;
            console.log(this.selectedNodeData);
          },
          error: (err: any) => {
            this.hf_service.handleError(err);
          },
          complete: () => {
            this.node_loading = false;
            if(!this.isLoading()) {
              this.data = this.processData();
              this.displayData = this.data;
              this.genFilterData();
            }
            console.log('Got Node_Info');
          }
        });
      }
    });

    this.hf_service.getHotfixData().subscribe({
      next: (rdata: CloudHotfixManifest) => {
        this.hotfix_info = rdata;
        console.log(this.hotfix_info);
      },
      error: (err: any) => {
        this.hf_service.handleError(err);
      },
      complete: () => {
        this.hotfix_loading = false;
        if(!this.isLoading()) {
          this.data = this.processData();
          this.displayData = this.data;
          this.genFilterData();
        }
        console.log('Got Hotfix_Info');
      }
    });
      
  }

  isLoading(): boolean {
    return (this.node_loading || this.hotfix_loading);
  }

  sortData(colName: string): void {
    if(this.sortColumn != colName) {
      this.sortColumn = colName;
      this.sortOrder = "";
    }
    if (this.sortOrder === "" || this.sortOrder === "dsc"){
      this.sortOrder = "asc";
      switch(colName) {
        case 'impactedArea':
          this.displayData.sort((a: any,b: any) => a[colName].join(', ') > b[colName].join(', ') ? 1 : -1);
          break;
        case 'requiredActions':
          this.displayData.sort((a: any,b: any) => this.getRequiredActionString(a[colName]) > this.getRequiredActionString(b[colName]) ? 1 : -1);
          break;
        default:
          this.displayData.sort((a: any,b: any) => a[colName] > b[colName] ? 1 : -1);
          break;
      }
      
    } else if (this.sortOrder === "asc"){
      this.sortOrder = "dsc";
      switch(colName) {
        case 'impactedArea':
          this.displayData.sort((a: any,b: any) => a[colName].join(', ') < b[colName].join(', ') ? 1 : -1);
          break;
        case 'requiredActions':
          this.displayData.sort((a: any,b: any) => this.getRequiredActionString(a[colName]) < this.getRequiredActionString(b[colName]) ? 1 : -1);
          break;
        default:
          this.displayData.sort((a: any,b: any) => a[colName] < b[colName] ? 1 : -1);
          break;
      }
    }
  }

  genFilterData():void {
    this.filterColumns.forEach((ele: string) => {
      var temp_data: Array<string> = new Array<string>();
      this.data.forEach((val:any) => {
        let t: string[]= [];
        if (ele === 'impactedArea') {
          t = val[ele];
        } else if (ele === 'requiredActions') {
          if(val[ele]['systemReboot'] == 'Yes') {
            t.push('System Reboot');
          } else if (val[ele]['productRestart'] == 'Yes'){
            t.push('Product Restart');
          } else if (val[ele]['serviceRestart'].length > 0){
            for(let i=0; i<val[ele]['serviceRestart'].length; i++){
              t.push('Service Restart: '+ val[ele]['serviceRestart'][i])
            }
          }
        } else {
          t.push(val[ele]);
        }
        temp_data = [...temp_data, ...t]
      });
      //console.log('New Filter Data: ', temp_data);
      let unique_temp_data = temp_data.filter((val, pos) => temp_data.indexOf(val) === pos);
      let filter_val = new Map<string, boolean>();
      unique_temp_data.forEach(val => filter_val.set(val, false));
      this.filterData.set(ele, filter_val);
    });
    //console.log(this.filterData);
  }

  onFilterChange(filterCol: string, filterVal: any, $event: any): void {
    this.filterApplied = false;
    this.collapseAll();
    this.clearSearch();
    this.filterData.get(filterCol)?.set(filterVal, $event.currentTarget.checked);
    this.displayData = this.data.filter((ele: any) => {
      var selected = true;
      this.filterData.forEach((fdval: any, fdkey: string) => {
        var tmpSelected = false;
        var filterPresent = false;
        fdval.forEach((val: boolean, key: string) => {
          if(val) {
            filterPresent = true;
            switch(fdkey) {
              case 'impactedArea': 
                tmpSelected = tmpSelected || ele[fdkey].includes(key);
                break;
              case 'requiredActions':
                switch(key) {
                  case 'System Reboot':
                    tmpSelected = tmpSelected || ele[fdkey]['systemReboot'] === 'Yes';
                    break;
                  case 'Product Restart':
                    tmpSelected = tmpSelected || ele[fdkey]['productRestart'] === 'Yes';
                    break;
                  default :
                    let kv = key.split(': ');
                    tmpSelected = tmpSelected || ele[fdkey]['serviceRestart'].includes(kv[1]);
                    break;
                }
                break;
              default:
                tmpSelected = tmpSelected || ele[fdkey] === key;
                break;
            }
          } 
        });
        if (filterPresent) {
          this.filterApplied = true;
          selected = selected && tmpSelected;
        }
      });
      return selected;
    });
  }

  clearFilter(): void {
    this.filterData.forEach((value: Map<string, boolean>, key: string) => {
      value.forEach((_, k: string) => this.filterData.get(key)?.set(k, false));
    });
    this.filterApplied = false;
  }

  clearSearch(): void {
    this.search.nativeElement.value = '';
    this.searchApplied = false;
  }

  clearSort(): void {
    this.sortData('idx');
    this.sortColumn = '';
    this.sortOrder = '';
  }
  clearAll(): void {
    this.collapseAll();
    this.displayData = this.data;
    if(this.filterApplied){
      this.clearFilter();
    }
    if(this.searchApplied){
      this.clearSearch()
    }
    if(this.sortColumn != ''){
      this.clearSort();
    }
  }

  expand(ele: any, id:number): void {
    this.displayData[id].show = !ele.currentTarget.classList.contains('collapsed')
  }

  collapseAll(): void {
    this.displayData.forEach((ele: any) => {
      if(ele.show){
        ele.show = false;
        document.getElementById('collapse'+ele.idx)?.classList.remove('show');
      }
    });
  }

  onSearch(search: any): void {
    var searchStr: string = search.currentTarget.value;
    searchStr = searchStr.trim().toLowerCase();
    if(searchStr.length == 0) {
      this.searchApplied = false;
      this.displayData = this.data;
      return;
    }
    this.collapseAll();
    this.clearFilter();
    this.searchApplied = true;
    this.displayData = this.data.filter((ele: any) => {
      for(let k in ele){
        if(typeof ele[k] != 'boolean' && ele[k].toString().toLowerCase().match(searchStr)) {
          return true;
        }
      }
      return false;
    });
  }

  findNode(search: any): void {
    var searchStr: string = search.currentTarget.value;
    searchStr = searchStr.trim().toLowerCase();
    if(searchStr.length == 0) {
      this.listNodes = this.allNodes;
      return;
    }
    this.listNodes = this.allNodes.filter((ele: any) => {
      for(let k in ele){
        if(typeof ele[k] != 'boolean' && ele[k].toString().toLowerCase().match(searchStr)) {
          return true;
        }
      }
      return false;
    });
  }

  selectNode(n: NodeData) {
    if(this.selectedNode !== n) {
      this.nodeListMenu.nativeElement.classList.remove('show');
      this.node_loading = true;
      this.selectedNode = n;
      this.searchNode.nativeElement.value = '';
      this.listNodes = this.allNodes;
      this.clearAll();
      this.hf_service.getNodeInfo(this.selectedNode.unique_id).subscribe({
        next: (rdata: NodeHotfixData) => {
          this.selectedNodeData = rdata;
          console.log(this.selectedNodeData);
        },
        error: (err: any) => {
          this.hf_service.handleError(err);
        },
        complete: () => {
          this.node_loading = false;
          this.data = this.processData();
          this.displayData = this.data;
          this.genFilterData();
          console.log('Got Node_Info');
        }
      });
    }
  }

  processData() {
    let temp_data: ProcessedHotfix[] = [];
    this.hotfixSummary['IMPORTANT']['available'] = 0;
    this.hotfixSummary['IMPORTANT']['installed'] = 0;
    this.hotfixSummary['RECOMMENDED']['available'] = 0;
    this.hotfixSummary['RECOMMENDED']['installed'] = 0;
    this.hotfixSummary['OPTIONAL']['available'] = 0;
    this.hotfixSummary['OPTIONAL']['installed'] = 0;

    this.hotfix_info.data.forEach((hf_data: Hotfix, index: number) => {

      if(hf_data.compatibleNode === 'ALL' || hf_data.compatibleNode === this.selectedNodeData.role) {
        let temp_hf: ProcessedHotfix;
        // temp_hf = Object.assign(hf_data);
        temp_hf = {...hf_data, applyStatus: 'Not Available', applyTimestamp: '', idx: index, show: false};
        // temp_hf.applyStatus = 'Not Available';
        // temp_hf.applyTimestamp = '';
        this.hotfixSummary[hf_data.severity]['available'] += 1;
        if(this.selectedNodeData.status === 'ONLINE') {
          temp_hf.applyStatus = 'Not Installed';
          this.selectedNodeData.hotfixes.forEach((hf: any) => {
            if(hf.status == 'SUCCESS') {
              if (hf_data.name === hf.name) {
                if (temp_hf.applyTimestamp < hf.timestamp) {
                  temp_hf.applyStatus = 'Installed';
                  temp_hf.applyTimestamp = hf.timestamp;
                }
                this.hotfixSummary[hf_data.severity]['installed'] += 1;
              } else if (hf_data.revert && hf_data.revert.name === hf.name) {
                if (temp_hf.applyTimestamp < hf.timestamp) {
                  temp_hf.applyStatus = 'Reverted';
                  temp_hf.applyTimestamp = hf.timestamp;
                }
                this.hotfixSummary[hf_data.severity]['installed'] -= 1;
              } 
            }
          });
        }
        temp_data.push(temp_hf);
        if(temp_hf.applyStatus == 'Not Installed' && temp_hf.severity != 'OPTIONAL') {
          this.sendNotification(temp_hf);
        }
      }
    });
    console.log(temp_data);
    return temp_data;
  }

  getRequiredActionString(action: HotfixRequiredActions): string {
    if(action.systemReboot.toLowerCase().match('yes')){
      return 'System Reboot';
    } else if (action.productRestart.toLowerCase().match('yes')) {
      return 'Product Restart';
    } else if (action.serviceRestart.length > 0) {
      return 'Service Restart(' + action['serviceRestart'].join(', ') + ')';
    }
    return '-';
  }

  getColorClass(applyStatus: string, severity: string): string {
    if(applyStatus === 'Installed'){
      return 'installed';
    }
    if(applyStatus === 'Not Installed'){
      if(severity === 'IMPORTANT'){
        return 'imp-not-installed';
      }
      if(severity === 'RECOMMENDED'){
        return 'rmd-not-installed';
      }
    }
    return 'default';
  }

  getFixesString(fixes: HotfixFixes): string {
    let result: string = 'General Enhancements';
    if(fixes.BUGFIX.length > 0){
      result = 'Bugfixes';
    }
    if(fixes.CVE.length > 0) {
      let tmp: string = 'CVE Fixes(';
      fixes.CVE.forEach((e: any, i: number) => {
        if(i != fixes.CVE.length-1){
          tmp = tmp + e.id + ', ';
        } else {
          tmp = tmp + e.id + ')';
        }
      });
      result = result + ', ' + tmp;
    }
    //TODO: is security format is same as CVE format then replicate the logic
    if(fixes.SECURITY.length > 0) {
      result = result + ', ' + 'Security Fixes';
    }
    return result;
  }

  sendNotification(hf: Hotfix) {
    let msg: string;
    switch(hf.severity) {
      case 'RECOMMENDED':
        msg = "A recommended hotfix is available, please apply the hotfix to improve system stablity.";
        break;
      case 'IMPORTANT':
        msg = "A important hotfix is available, it is recommended to apply the hotfix as soon as possible.";
        break;
      default:
        msg = "A hotfix is available, please apply the hotfix to improve system stablity.";
        break;
    }
    let n: Notification = {
      source: "HOTFIX", 
      type: hf.severity,
      tag: this.selectedNodeData.hostname, 
      title: "Hotfix-Alert", 
      body: msg, 
      action: "DEFAULT", 
      extraData: hf
    };
    this.notify_service.sendNotification(n);
  }

  testFunc(): void {

  }
}
