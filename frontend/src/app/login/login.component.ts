import { Component, OnInit } from '@angular/core';
import * as OktaSignIn from '@okta/okta-signin-widget';
import oktaConfig from '../okta.config';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnInit {
  signIn: any;

  constructor() {
    this.signIn = new OktaSignIn({
      baseUrl: oktaConfig.issuer.split('/oauth2')[0],
      clientId: oktaConfig.clientId,
      redirectUri: oktaConfig.redirectUri,
      // logo: '/assets/go.png',
      i18n: {
        en: {
          'primaryauth.title': 'Sign in to Go Blog App',
        },
      },
      authParams: {
        responseType: ['id_token', 'token'],
        issuer: oktaConfig.issuer,
        display: 'page',
        scopes: oktaConfig.scope.split(' '),
      },
    });
  }

  ngOnInit() {
    this.signIn.renderEl(
      { el: '#sign-in-widget' },
      () => {
        /**
         * In this flow, the success handler will not be called because we redirect
         * to the Okta org for the authentication workflow.
         */
      },
      (err) => {
        throw err;
      },
    );
  }

}
