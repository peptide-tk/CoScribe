import React from "react";
import { Header } from "./Header";
import { Footer } from "./Footer";

interface LayoutProps {
  children: React.ReactNode;
}

export const Layout: React.FC<LayoutProps> = ({ children }) => {
  return (
    <div className="App">
      <Header />
      <main className="app-main">
        {children}
      </main>
      <Footer />
    </div>
  );
};