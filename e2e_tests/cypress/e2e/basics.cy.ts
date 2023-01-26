/// <reference types="cypress" />

const BASE_URL = 'http://localhost:8080';
const CLEOC_PATH = '../dist/cleoc_linux_amd64_v1/cleoc';

describe('basic front page tests', () => {
  beforeEach(() => {
    cy.visit(BASE_URL)
  });

  it('passes', () => {
    cy.findByLabelText('Title *').click().type('Does cypress work?')
  });

  it('renders the page title and header', () => {
    cy.title().should('eq', 'Cleodora');
    cy.findAllByRole('heading').first().should('have.text', 'Cleodora');
  });

  it('contains new forecast', () => {
    const cmd = CLEOC_PATH + " add forecast " +
        "-t 'Is this a test forecast?' " +
        "-r '2022-12-01T15:00:00+01:00' " +
        "--reason \"We're running a test, so it seems likely.\" " +
        "-p Yes=99 " +
        "-p No=1";
    cy.exec(cmd)
        .its('stdout')
        .should('not.be.empty');
    cy.reload();
    cy.findByText("Is this a test forecast?");
  });
})
