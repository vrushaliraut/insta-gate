import { render, screen } from "@testing-library/react";
import Home from "./page";
// Some linters prefer you name it Home if the default export is Home()
describe("Home Page", () => {
  it("renders successfully", () => {
    render(<Home />);
    // Next.js 14+ default page usually contains a main element
    const mainElement = screen.getByRole("main");
    expect(mainElement).toBeInTheDocument();
  });
});
