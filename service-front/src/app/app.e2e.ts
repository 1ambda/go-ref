import { browser, by, element } from 'protractor';
import 'tslib';

describe('App', () => {

  beforeEach(async () => {
    await browser.get('/');
  });

  // it('should have header', async () => {
  //   let subject = await element(by.css('h1')).isPresent();
  //   let result  = true;
  //   expect(subject).toEqual(result);
  // });
});
