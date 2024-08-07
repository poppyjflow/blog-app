import React from 'react';
import { Link } from 'react-router-dom';
import styled from 'styled-components';

const HeaderWrapper = styled.header`
  background-color: ${props => props.theme.primary};
  padding: 1rem;
`;

const Nav = styled.nav`
  display: flex;
  justify-content: space-between;
`;

const NavLink = styled(Link)`
  color: ${props => props.theme.text};
  text-decoration: none;
  margin-right: 1rem;
`;

function Header() {
  return (
    <HeaderWrapper>
      <Nav>
        <NavLink to="/">Home</NavLink>
        <NavLink to="/posts">Posts</NavLink>
        <NavLink to="/create">Create Post</NavLink>
      </Nav>
    </HeaderWrapper>
  );
}

export default Header;