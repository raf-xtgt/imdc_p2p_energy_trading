import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ManageClerksComponent } from './manage-clerks.component';

describe('ManageClerksComponent', () => {
  let component: ManageClerksComponent;
  let fixture: ComponentFixture<ManageClerksComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ManageClerksComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ManageClerksComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
