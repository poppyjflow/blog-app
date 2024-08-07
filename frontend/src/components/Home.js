import React from 'react';
import styled from 'styled-components';

const HomeWrapper = styled.div`
  padding: 2rem;
`;

function Home() {
  return (
    <HomeWrapper>
      <h1>Welcome to the Blog App</h1>
      <p>Start exploring posts or create your own!</p>
    </HomeWrapper>
  );
}

export default Home;