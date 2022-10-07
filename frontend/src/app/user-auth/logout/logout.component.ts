import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-logout',
  templateUrl: './logout.component.html',
  styleUrls: ['./logout.component.css']
})
export class LogoutComponent implements OnInit {

  session_timeout: Boolean = false;
  constructor(private activatedRoute: ActivatedRoute) { }

  ngOnInit(): void {
    this.session_timeout = Boolean(this.activatedRoute.snapshot.queryParamMap.get('session_timeout')) || false;
    console.log("Logout called:", this.session_timeout);
  }

}
