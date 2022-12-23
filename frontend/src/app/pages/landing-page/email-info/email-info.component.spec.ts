import { ComponentFixture, TestBed } from '@angular/core/testing';

import { EmailInfoComponent } from './email-info.component';

describe('EmailInfoComponent', () => {
  let component: EmailInfoComponent;
  let fixture: ComponentFixture<EmailInfoComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ EmailInfoComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(EmailInfoComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
