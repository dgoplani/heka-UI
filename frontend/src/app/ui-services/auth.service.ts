import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from 'src/environments/environment';

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  private server: string;
  
  constructor(private http : HttpClient) { 
    if(environment.production){
      this.server = "https://"+window.location.host+"/bloxconnect/api";
    } else {
      this.server = "http://127.0.0.1:5200";
    }
  }

  login(username: string, password: string) {
    var body = {
      "username": username, 
      "password": password
    };
    var req_headers = new HttpHeaders();
    req_headers.append("Content-type", 'application/json');
    return this.http.post(this.server+"/login", body, {headers:req_headers, observe: 'response', withCredentials: true});
  }

  logout() {
    var body = undefined;
    var req_headers = new HttpHeaders();
    req_headers.append("Content-type", 'application/json');
    return this.http.post(this.server+"/logout", body, {headers:req_headers, observe: 'response', withCredentials: true});
  }

  setToken(username: string) {
    localStorage.setItem('bloxconnect_loggedin', username);
  }

  removeToken() {
    localStorage.clear();
  }

  isAuthenticated() {
    return !!localStorage.getItem('bloxconnect_loggedin');
  }

}
