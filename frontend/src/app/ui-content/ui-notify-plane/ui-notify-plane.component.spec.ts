import { ComponentFixture, TestBed } from '@angular/core/testing';

import { UiNotifyPlaneComponent } from './ui-notify-plane.component';

describe('UiNotifyPlaneComponent', () => {
  let component: UiNotifyPlaneComponent;
  let fixture: ComponentFixture<UiNotifyPlaneComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ UiNotifyPlaneComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(UiNotifyPlaneComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
