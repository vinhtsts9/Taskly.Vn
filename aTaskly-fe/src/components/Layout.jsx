import React from "react";
import { Outlet } from "react-router-dom";
import Header from "./Header";
import Footer from "./Footer";
import AIChatbot from "./AIChatbot";

const Layout = () => {
  return (
    <>
      <Header />
      <main>
        <Outlet />
      </main>
      <Footer />
      <AIChatbot />
    </>
  );
};

export default Layout;
