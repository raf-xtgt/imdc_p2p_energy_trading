import { ComponentFixture, TestBed } from '@angular/core/testing';

import { NewMarketPageComponent } from './new-market-page.component';

describe('NewMarketPageComponent', () => {
  let component: NewMarketPageComponent;
  let fixture: ComponentFixture<NewMarketPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ NewMarketPageComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(NewMarketPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
