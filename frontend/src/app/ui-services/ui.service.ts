import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from 'src/environments/environment';
import { GridData } from '../interfaces/grid-data';
import { NodeData } from '../interfaces/node-data';

@Injectable({
  providedIn: 'root'
})
export class UiService {

  private server: string;

  constructor(private http : HttpClient) {
    if(environment.production){
      this.server = "https://"+window.location.host+"/bloxconnect/api";
    } else {
      this.server = "http://127.0.0.1:5200";
    }
  }

  getGridInfo(): Observable<GridData> {
    var req_headers = new HttpHeaders();
    req_headers.append("Content-type", 'application/json');
    return this.http.get<GridData>(this.server+"/grid_data", {headers:req_headers, observe: 'body', withCredentials: true});
  }

  getNodeList(): Observable<NodeData[]> {
    var req_headers = new HttpHeaders();
    req_headers.append("Content-type", 'application/json');
    return this.http.get<NodeData[]>(this.server+"/list_nodes", {headers:req_headers, observe: 'body', withCredentials: true});
  }
}
