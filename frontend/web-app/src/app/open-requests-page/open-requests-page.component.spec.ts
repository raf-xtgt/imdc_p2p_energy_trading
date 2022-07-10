import { ComponentFixture, TestBed } from '@angular/core/testing';

import { OpenRequestsPageComponent } from './open-requests-page.component';

describe('OpenRequestsPageComponent', () => {
  let component: OpenRequestsPageComponent;
  let fixture: ComponentFixture<OpenRequestsPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ OpenRequestsPageComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(OpenRequestsPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
