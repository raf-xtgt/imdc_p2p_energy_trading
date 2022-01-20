import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MarketpageComponent } from './marketpage.component';

describe('MarketpageComponent', () => {
  let component: MarketpageComponent;
  let fixture: ComponentFixture<MarketpageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ MarketpageComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(MarketpageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
