import React, { useState } from 'react';
import axios from 'axios';
import styled from 'styled-components';

const Form = styled.form`
  display: flex;
  flex-direction: column;
  max-width: 500px;
  margin: 2rem auto;
`;

const Input = styled.input`
  margin-bottom: 1rem;
  padding: 0.5rem;
  background-color: ${props => props.theme.background};
  color: ${props => props.theme.text};
  border: 1px solid ${props => props.theme.text};
`;

const TextArea = styled.textarea`
  margin-bottom: 1rem;
  padding: 0.5rem;
  background-color: ${props => props.theme.background};
  color: ${props => props.theme.text};
  border: 1px solid ${props => props.theme.text};
`;

const Button = styled.button`
  background-color: ${props => props.theme.primary};
  color: ${props => props.theme.text};
  padding: 0.5rem;
  border: none;
  cursor: pointer;
`;

function CreatePost() {
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const user2 = 2;

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await axios.post('http://localhost:8080/posts', { user2, title, content });
      setTitle('');
      setContent('');
      alert('Post created successfully!');
    } catch (error) {
      console.error('Error creating post:', error);
    }
  };

  return (
    <Form onSubmit={handleSubmit}>
      <h2>Create a New Post</h2>
      <Input
        type="text"
        placeholder="Title"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        required
      />
      <TextArea
        placeholder="Content"
        value={content}
        onChange={(e) => setContent(e.target.value)}
        required
      />
      <Button type="submit">Create Post</Button>
    </Form>
  );
}

export default CreatePost;