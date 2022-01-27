import { ComponentFixture, TestBed } from '@angular/core/testing';

import { HouseholdPriceGraphComponent } from './household-price-graph.component';

describe('HouseholdPriceGraphComponent', () => {
  let component: HouseholdPriceGraphComponent;
  let fixture: ComponentFixture<HouseholdPriceGraphComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ HouseholdPriceGraphComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(HouseholdPriceGraphComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
