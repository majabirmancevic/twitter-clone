import { ComponentFixture, TestBed } from '@angular/core/testing';

import { BusinessProfilePageComponent } from './business-profile-page.component';

describe('BusinessProfilePageComponent', () => {
  let component: BusinessProfilePageComponent;
  let fixture: ComponentFixture<BusinessProfilePageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ BusinessProfilePageComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(BusinessProfilePageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
