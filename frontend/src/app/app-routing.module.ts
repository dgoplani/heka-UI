import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { LoginComponent } from './user-auth/login/login.component';
import { UiComponent } from './ui-content/ui/ui.component';
import { NotFoundComponent } from './not-found/not-found.component';
import { LogoutComponent } from './user-auth/logout/logout.component';
import { UiHotfixComponent } from './ui-content/ui-hotfix/ui-hotfix.component';
import { AuthGuard } from './ui-guards/auth.guard';

const routes: Routes = [
  { path: '', 
    redirectTo:'ui', 
    pathMatch: 'full'
  },
  {  
    path: 'login', 
    component: LoginComponent,
    data: {
      title: 'BloxConnect - Login'
    }
  },
  {
    path: 'logout', 
    component: LogoutComponent,
    data: {
      title: 'BloxConnect - Logout'
    }
  },
  {
    path: 'ui', 
    component: UiComponent,
    canActivate: [AuthGuard],
    children: [
      {
        path: '', 
        redirectTo:'hotfix', 
        pathMatch:'full'
      },
      {
        path: 'hotfix', 
        component: UiHotfixComponent,
        data: {
          title: 'BloxConnect - Hotfix Information'
        }
      }
    ]
  },
  {
    path: '**', 
    component: NotFoundComponent,
    data: {
      title: 'BloxConnect - Not Found'
    }
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
