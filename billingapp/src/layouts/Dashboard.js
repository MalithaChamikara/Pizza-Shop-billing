import React from "react";
import { Grid2 } from "@mui/material";
import Header from "../components/Header.js";
import SideBar from "../components/SideBar.js";

const Dashboard = ({ children }) => {
  return (
    <div>
      <Header />
      <Grid2 container sx={{ mt: 2, mb: 2 }}>
        <Grid2 item xs={2}>
          <SideBar />
        </Grid2>
        <Grid2 item xs={10}>
          {children}
        </Grid2>
      </Grid2>
    </div>
  );
};

export default Dashboard;
