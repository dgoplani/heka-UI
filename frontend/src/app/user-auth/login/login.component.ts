import { Component, OnInit } from '@angular/core';
import { FormGroup, Validators, FormBuilder } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { AuthService } from 'src/app/ui-services/auth.service';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {

  login_failed: boolean = false;
  authenticating: boolean = false;
  redirectUrl: any = '';
  loginForm!: FormGroup;

  constructor(private formBuilder: FormBuilder, 
              private auth: AuthService, 
              private router: Router,
              private active_route: ActivatedRoute) {
    if(this.auth.isAuthenticated()){
      this.router.navigate(['ui']);
    }
  }

  ngOnInit(): void {
    this.redirectUrl = this.active_route.snapshot.queryParamMap.get('redirectUrl') || '/';
    this.loginForm = this.formBuilder.group({
      username: ['', Validators.required],
      password: ['', Validators.required]
    })
  }

  onSubmit(): void {
    console.log("Login called");
    this.login_failed = false;
    if(this.loginForm.value.username === "" || this.loginForm.value.password === "") {
      this.login_failed = true;
      return;
    }
    this.authenticating = true;
    this.auth.login(this.loginForm.value.username, this.loginForm.value.password).subscribe({
      next: (response: any) => {
        console.log(response);
        //console.log(response.headers);
      },
      error: (err: any) => {
        console.log(err);
        this.authenticating = false;
        this.login_failed = true;
      },
      complete: () => {
        this.authenticating = false;
        this.login_failed = false;
        this.auth.setToken(this.loginForm.value.username);
        this.router.navigateByUrl(this.redirectUrl);
      }
    });
  }
}