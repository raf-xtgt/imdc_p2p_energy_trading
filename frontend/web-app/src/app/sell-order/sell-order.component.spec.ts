import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SellOrderComponent } from './sell-order.component';

describe('SellOrderComponent', () => {
  let component: SellOrderComponent;
  let fixture: ComponentFixture<SellOrderComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ SellOrderComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(SellOrderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
