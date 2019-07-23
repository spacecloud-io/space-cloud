import React from 'react';
import Routes from './Routes';
import discord from './assets/discord.svg'

function App() {
  return (
    <div>
      <Routes />
      <a href="https://discordapp.com/invite/ypXEEBr" target="_blank" class="discord valign-wrapper">
        <span>Have a Question?</span>
        <img src={discord} alt="" />
      </a>
    </div>
  );
}

export default App;
