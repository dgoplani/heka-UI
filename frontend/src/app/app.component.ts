import { Component, OnInit } from '@angular/core';
import { environment } from 'src/environments/environment';
import { Title } from '@angular/platform-browser';
import { ActivatedRoute, NavigationEnd, Router } from '@angular/router';
import { filter } from 'rxjs';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit{

  constructor(private router: Router,  
              private activatedRoute: ActivatedRoute,  
              private titleService: Title) {
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
        console.log(data);  
        this.titleService.setTitle(data['title'])});  
    });  
  }

  getChild(activatedRoute: ActivatedRoute): ActivatedRoute {  
    if (activatedRoute.firstChild) {  
      return this.getChild(activatedRoute.firstChild);  
    } else {  
      return activatedRoute;  
    }  
  }  
}
