import { Component, OnInit } from '@angular/core';
import { environment } from 'src/environments/environment';
import { Title } from '@angular/platform-browser';
import { ActivatedRoute, NavigationEnd, Router } from '@angular/router';
import { filter, interval, Subscription } from 'rxjs';
import { BnNgIdleService } from 'bn-ng-idle';
import { AuthService } from './ui-services/auth.service';
import { UiService } from './ui-services/ui.service';
import { ReadyStatus } from './interfaces/ready-status';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit{

  isReady: boolean = false;
  readySub!: Subscription;
  constructor(private router: Router,  
              private activatedRoute: ActivatedRoute,  
              private titleService: Title,
              private auth: AuthService,
              private uiService: UiService,
              private userIdleService: BnNgIdleService) {
    if(environment.production){
      console.log('Production Mode');
    } else {
      console.log('Development Mode');
    }
  }

  ngOnInit() {
    this.router.events.pipe(  
      filter(event => event instanceof NavigationEnd),  
    ).subscribe(() => {  
      const rt = this.getChild(this.activatedRoute);  
      rt.data.subscribe(data => {  
        //console.log(data);  
        this.titleService.setTitle(data['title'])});  
    });  

    // TODO : Make session timeout sync from NIOS instead of hardcoded 10 mins
    this.userIdleService.startWatching(600).subscribe((isTimedOut: Boolean) => {
      //console.log("Timeout: ", isTimedOut);
      if(this.auth.isAuthenticated() && isTimedOut) {
        console.log("Session Expired: Logging Out");
        this.auth.logout().subscribe({
          next: (response: any) => {
            
          },
          error: (err: any) => {
            console.log(err);
            localStorage.clear();
            this.router.navigate(['/login'])
          },
          complete: () => {
            this.router.navigate(['/logout'], { queryParams: { session_timeout: true } });
            this.auth.removeToken();
          }
        })
      } else {
        console.log("User Not Logged In.");
      }
    });

    this.uiService.getReadyStatus().subscribe((res: ReadyStatus) => {
      if(res.ready){
        this.isReady = res.ready;
      }
    });
    if(this.isReady == false){
      this.readySub = interval(10000).subscribe(n => {
        this.uiService.getReadyStatus().subscribe((res: ReadyStatus) => {
          if(res.ready){
            this.isReady = res.ready;
            this.readySub.unsubscribe();
          } else {
            console.log("UI is not Ready, Try", n)
          }
        });
      });
    }
  }

  getChild(activatedRoute: ActivatedRoute): ActivatedRoute {  
    if (activatedRoute.firstChild) {  
      return this.getChild(activatedRoute.firstChild);  
    } else {  
      return activatedRoute;  
    }  
  }  
}
