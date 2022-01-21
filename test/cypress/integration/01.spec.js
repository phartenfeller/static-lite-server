const URL = "http://localhost:3002";

describe("My First Test", () => {
  it("Visit index page", () => {
    cy.visit(URL);

    cy.get("body").then(($body) => {
      const $header = $body.find("#header");
      expect($header).to.contain("Hello World");
    });

    cy.request("/").then((response) => {
      expect(response.status).to.equal(200);
      expect(response).to.have.property("headers");
      cy.log("headers", response.headers);

      // headers are lowercase
      expect(response.headers).to.have.property("server");
      expect(response.headers["server"]).to.contain("static-lite-server:");
    });
  });
});
