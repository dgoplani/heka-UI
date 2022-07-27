import { Injectable } from '@angular/core';
import { Subject } from 'rxjs';
import { Notification } from '../interfaces/notification';

@Injectable({
  providedIn: 'root'
})
export class NotificationService {

  notifications!: Subject<Notification>;
  constructor() {
    this.notifications = new Subject<Notification>();
  }

  sendNotification(ndata: Notification) {
    this.notifications.next(ndata);
  }

}
