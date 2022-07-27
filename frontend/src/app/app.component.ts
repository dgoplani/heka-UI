import { Component } from '@angular/core';
import { environment } from 'src/environments/environment';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {

  title = 'BloxConnect';

  constructor() {
    if(environment.production){
      console.log('Production Mode');
    } else {
      console.log('Development Mode');
    }
  }

}
