describe('resolve forecast', () => {
  let cleocPath:string;

  before(() => {
    cleocPath = Cypress.env('cleocPath') as string
    if (cleocPath == undefined) {
      throw new Error("cleocPath is not defined in Cypress.env");
    }
  });

  beforeEach(() => {
    cy.visit("/");
  });

  it('succeeds', () => {
    // Add a new forecast
    // The purpose of randomNum is to allow running the test multiple times
    // since every forecast can only be resolved once.
    const randomNum = Math.floor(Math.random() * 10000);
    const title = `Will this forecast be resolved? (${randomNum})`;
    const resolves = new Date();
    resolves.setDate(resolves.getDate() + 1);
    const resolvesISO = resolves.toISOString();
    const cmd = cleocPath + " " +
        "--url " + Cypress.config('baseUrl') + " " +
        "add forecast " +
        `-t '${title}' ` +
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

    // Ensure no dialog is open
    cy.findByRole('dialog').should('not.exist');

    // Select the correct row in the table of forecasts
    cy.findByText(title).parent("tr").as('row');

    // Ensure the forecast is not yet resolved and open the resolve dialog
    cy.get('@row').within(() => {
      cy.findByText('UNRESOLVED').should('exist');
      cy.findByRole('button', { name: /Resolve/i }).click();
    });

    // Ensure the resolve dialog contains the right stuff.
    cy.findByRole('dialog').within(() => {
      cy.findByLabelText('Yes').should('exist');
      cy.findByLabelText('No').should('exist');
      cy.findByLabelText('No correct outcome (resolve forecast as N/A)').should('exist');
    });

    // Close the resolve dialog and ensure it's gone.
    cy.findByRole('button', { name: "Cancel" }).click();
    cy.findByRole('dialog').should('not.exist');

    // Open the dialog again.
    cy.get('@row').within(() => {
      cy.findByRole('button', { name: /Resolve/i }).click();
    });

    // Choose the outcome 'No' and save.
    cy.findByRole('dialog').within(() => {
      cy.findByLabelText('No').should('exist').click();
      cy.findByRole('button', { name: 'Save' }).click();
    });

    // Verify the resolution is RESOLVED and the correct outcome is bold.
    cy.get('@row').within(() => {
      cy.findByText('RESOLVED').should('exist');
      cy.findByText('No: 1%').should('have.css', 'font-weight', '700');
    });
  })
})