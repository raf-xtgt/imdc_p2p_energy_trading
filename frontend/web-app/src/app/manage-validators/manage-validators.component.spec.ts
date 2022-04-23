import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ManageValidatorsComponent } from './manage-validators.component';

describe('ManageValidatorsComponent', () => {
  let component: ManageValidatorsComponent;
  let fixture: ComponentFixture<ManageValidatorsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ManageValidatorsComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ManageValidatorsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
