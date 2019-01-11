import { TestBed } from '@angular/core/testing';

import { Auth.InterceptorService } from './auth.interceptor.service';

describe('Auth.InterceptorService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: Auth.InterceptorService = TestBed.get(Auth.InterceptorService);
    expect(service).toBeTruthy();
  });
});
