import { Component, OnInit } from '@angular/core';
import { UiService } from 'src/app/ui-services/ui.service';
import { Router, ActivatedRoute } from '@angular/router';
import { UiHotfixComponent } from '../ui-hotfix/ui-hotfix.component';
import { GridData } from 'src/app/interfaces/grid-data';
import { NodeData } from 'src/app/interfaces/node-data';

@Component({
  selector: 'app-ui',
  templateUrl: './ui.component.html',
  styleUrls: ['./ui.component.css']
})

export class UiComponent implements OnInit {
 
  loading_grid_info: boolean = true;
  loading_node_list: boolean = true;
  loading_title: string = "Please wait, while we are getting things ready for you..."

  grid_info: any;
  node_list: any;
  hotfix_comp!: UiHotfixComponent;
  constructor(private ui_service: UiService, 
              private router: Router,
              private route: ActivatedRoute) {
  }

  ngOnInit(): void { 
    this.ui_service.getGridInfo().subscribe({
      next: (data: GridData) => {
        console.log(data);
        //console.log(response.headers);
        this.grid_info = data;
      },
      error: (err: any) => {
        console.log(err);
        localStorage.clear();
        this.router.navigate(['/login']);
      },
      complete: () => {
        this.loading_grid_info = false;
        console.log('Got Grid_info');
      }
    });

    this.ui_service.getNodeList().subscribe({
      next: (data: NodeData[]) => {
        console.log(data);
        //console.log(response.headers);
        this.node_list = data;
      },
      error: (err: any) => {
        console.log(err);
        localStorage.clear();
        this.router.navigate(['/login']);
      },
      complete: () => {
        this.loading_node_list = false;
        console.log('Got Node_List');
      }
    });
  }

  isLoading(): boolean {
    return this.loading_grid_info && this.loading_node_list;
  }

}
