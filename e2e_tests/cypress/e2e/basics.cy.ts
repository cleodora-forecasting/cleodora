/// <reference types="cypress" />

describe('basic front page tests', () => {
  let cleocPath:string;

  before(() => {
    cleocPath = Cypress.env('cleocPath') as string
    if (cleocPath == undefined) {
      throw new Error("cleocPath is not defined in Cypress.env");
    }
  });

  beforeEach(() => {
    cy.visit("/")
  });

  it('title input field for new forecast exists', () => {
    cy.findByLabelText('Title *').click().type('Does cypress work?')
  });

  it('renders the page title and header', () => {
    cy.title().should('eq', 'Cleodora');
    cy.findAllByRole('heading').first().should('have.text', 'Cleodora');
  });

  it('contains new forecast', () => {
    const resolves = new Date();
    resolves.setDate(resolves.getDate() + 1);
    const resolvesISO = resolves.toISOString();
    const cmd = cleocPath + " " +
        "--url " + Cypress.config('baseUrl') + " " +
        "add forecast " +
        "-t 'Is this a test forecast?' " +
        `-r '${resolvesISO}' ` +
        "--reason \"We're running a test, so it seems likely.\" " +
        "-p Yes=99 " +
        "-p No=1";
    cy.exec(cmd).then((result) => {
        expect(result.code).to.eq(0);
        expect(result.stderr).to.be.empty;
        expect(result.stdout).to.not.be.empty;
    });
    cy.reload();
    cy.findByText("Is this a test forecast?");
  });
})
