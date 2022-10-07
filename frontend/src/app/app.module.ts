import { NgModule } from '@angular/core';
import { BrowserModule, Title } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';

import { BnNgIdleService } from 'bn-ng-idle';

import { UserAuthModule} from './user-auth/user-auth.module'
import { UiContentModule } from './ui-content/ui-content.module';
import { NotFoundComponent } from './not-found/not-found.component';
import { HttpClientModule } from '@angular/common/http';
import { NotificationService } from './ui-services/notification.service';
import { AuthService } from './ui-services/auth.service';

@NgModule({
  declarations: [
    AppComponent,
    NotFoundComponent
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    AppRoutingModule,
    HttpClientModule,
    UserAuthModule,
    UiContentModule
  ],
  providers: [
    NotificationService, 
    AuthService,
    Title,
    BnNgIdleService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
