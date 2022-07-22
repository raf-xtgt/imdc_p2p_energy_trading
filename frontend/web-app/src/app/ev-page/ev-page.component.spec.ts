import { ComponentFixture, TestBed } from '@angular/core/testing';

import { EvPageComponent } from './ev-page.component';

describe('EvPageComponent', () => {
  let component: EvPageComponent;
  let fixture: ComponentFixture<EvPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ EvPageComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(EvPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
