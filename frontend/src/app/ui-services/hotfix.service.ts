import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { Observable } from 'rxjs';
import { environment } from 'src/environments/environment';
import { CloudHotfixManifest } from '../interfaces/cloud-hotfix-manifest';
import { NodeHotfixData } from '../interfaces/node-hotfix-data';

@Injectable({
  providedIn: 'root'
})
export class HotfixService {

  private server: string;

  constructor(private http : HttpClient, private router: Router) { 
    if(environment.production){
      //this.server = "https://"+window.location.host+"/bloxconnect/api";
      this.server = '/bloxconnect/api';
    } else {
      this.server = "http://127.0.0.1:5200";
    }
  }

  getHotfixData(): Observable<CloudHotfixManifest> {
    var req_headers = new HttpHeaders();
    req_headers.append("Content-type", 'application/json');
    return this.http.get<CloudHotfixManifest>(this.server+"/cloud_manifest", {headers:req_headers, observe: 'body', withCredentials: true});
  }

  getNodeInfo(id: string): Observable<NodeHotfixData> {
    var req_headers = new HttpHeaders();
    req_headers.append("Content-type", 'application/json');
    return this.http.get<NodeHotfixData>(this.server+"/node_data/"+id, {headers:req_headers, observe: 'body', withCredentials: true});
  }

  handleError(err: any) {
    if(err.status === 401) {
      localStorage.clear();
      this.router.navigate(['/login'])
    } else {
      console.log(err);
    }
  }
}


