import { Component, Input, OnInit, ViewChild } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from 'src/app/ui-services/auth.service';
import { UiNotifyPlaneComponent } from '../ui-notify-plane/ui-notify-plane.component';

@Component({
  selector: 'app-ui-header',
  templateUrl: './ui-header.component.html',
  styleUrls: ['./ui-header.component.css']
})
export class UiHeaderComponent implements OnInit {

  @Input('node_list')node_list!: any;
  @ViewChild(UiNotifyPlaneComponent, {
    static: false
  }) npanel!: UiNotifyPlaneComponent;   

  username!: string;

  constructor(private auth: AuthService, private router: Router) {  }

  ngOnInit(): void {
    this.username = localStorage.getItem('bloxconnect_loggedin') || '';
  }

  openNotifyPlane(): void {
    this.npanel.open();
  }

  logout() {
    this.auth.logout().subscribe({
      next: (response: any) => {
        
      },
      error: (err: any) => {
        console.log(err);
        localStorage.clear();
        this.router.navigate(['/login'])
      },
      complete: () => {
        this.auth.removeToken();
        this.router.navigate(['/logout']);
      }
    })
  }
}
