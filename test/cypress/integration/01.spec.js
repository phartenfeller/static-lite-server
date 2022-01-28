const URL = 'http://localhost:3002';
const NF_URL = URL + '/does-not-exists';

describe('My First Test', () => {
  it('Index shows correct content', () => {
    cy.visit(URL);

    cy.get('body').then(($body) => {
      const $header = $body.find('#header');
      expect($header).to.contain('Hello World');
    });
  });

  it('Check index status and header', () => {
    cy.visit(URL);
    cy.request('/').then((response) => {
      expect(response.status).to.equal(200);
      expect(response).to.have.property('headers');
      cy.log('headers', response.headers);

      // headers are lowercase
      expect(response.headers).to.have.property('server');
      expect(response.headers['server']).to.contain('static-lite-server:');
    });
  });

  it('Check default 404 page', () => {
    cy.visit(URL);
    cy.request({ url: NF_URL, failOnStatusCode: false }).then((response) => {
      expect(response.status).to.equal(404);

      expect(response.body).to.contain('Not Found (404)');
    });
  });
});
