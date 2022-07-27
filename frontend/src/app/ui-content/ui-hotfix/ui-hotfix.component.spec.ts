import { ComponentFixture, TestBed } from '@angular/core/testing';

import { UiHotfixComponent } from './ui-hotfix.component';

describe('UiBodyComponent', () => {
  let component: UiHotfixComponent;
  let fixture: ComponentFixture<UiHotfixComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ UiHotfixComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(UiHotfixComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
