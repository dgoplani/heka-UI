import { TestBed } from '@angular/core/testing';

import { HotfixService } from './hotfix.service';

describe('HotfixService', () => {
  let service: HotfixService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(HotfixService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
