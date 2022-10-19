import { NgModule } from '@angular/core';
import { CommonModule, DatePipe } from '@angular/common';
import { FaIconLibrary, FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { fas } from '@fortawesome/free-solid-svg-icons';
import { far } from '@fortawesome/free-regular-svg-icons';
import { NavbarComponent } from './navbar/navbar.component';
import { UiComponent } from './ui/ui.component';
import { UiHeaderComponent } from './ui-header/ui-header.component';
import { UiHotfixComponent } from './ui-hotfix/ui-hotfix.component';
import { UiNotifyPlaneComponent } from './ui-notify-plane/ui-notify-plane.component';
import { UiModalComponent } from './ui-modal/ui-modal.component';
import { AppRoutingModule } from '../app-routing.module';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { LoadingComponent } from '../loading/loading.component';
import { RouterModule } from '@angular/router';
import { UiNotificationComponent } from './ui-notification/ui-notification.component';

@NgModule({
  declarations: [
    NavbarComponent,
    UiComponent,
    UiHeaderComponent,
    UiHotfixComponent,
    UiNotifyPlaneComponent,
    LoadingComponent,
    UiModalComponent,
    UiNotificationComponent
  ],
  imports: [
    CommonModule,
    BrowserAnimationsModule,
    AppRoutingModule,
    FontAwesomeModule,
    RouterModule
  ],
  exports: [
    UiComponent
  ],
  providers: [
    DatePipe
  ]
})
export class UiContentModule { 
  constructor(lib: FaIconLibrary) {
    lib.addIconPacks(fas, far);
  }
}
