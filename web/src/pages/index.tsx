import {
  Button,
  Input,
  PressEvent,
  Spacer,
  Table,
  useInput,
} from '@nextui-org/react';
import { useState } from 'react';

interface BlogPost {
  author: string;
  date: string;
  title: string;
  url: string;
}

export default function Home() {
  const blogPostsColumns = [
    {
      key: 'title',
      label: 'Title',
    },
    {
      key: 'author',
      label: 'Author',
    },
    {
      key: 'date',
      label: 'Date',
    },
    {
      key: 'url',
      label: 'URL',
    },
  ];

  const blogPostsTableFC = (label: string, state: BlogPost[]) => {
    return (
      <Table
        aria-label={label}
      >
        <Table.Header columns={blogPostsColumns}>
          {(column) => (
            <Table.Column key={column.key}>{column.label}</Table.Column>
          )}
        </Table.Header>
        <Table.Body items={state}>
          {(item: BlogPost) => (
            <Table.Row key={item.url}>
              <Table.Cell>{item.title}</Table.Cell>
              <Table.Cell>{item.author}</Table.Cell>
              <Table.Cell>{item.date}</Table.Cell>
              <Table.Cell><a href={item.url}>{item.url}</a></Table.Cell>
            </Table.Row>
          )}
        </Table.Body>
      </Table>
    );
  };

  const [scrapedBlogPosts, setScrapedBlogPosts] = useState([]);

  const scrapeButtonOnPress = async (_: PressEvent) => {
    const resp = await fetch('http://localhost:8000/api/scrape');
    console.log("scrape response", resp);
    const data = await resp.json();
    setScrapedBlogPosts(data);
  };

  const newBlogPostTitle = useInput("");
  const newBlogPostAuthor = useInput("");
  const newBlogPostDate = useInput("");
  const newBlogPostURL = useInput("");

  const [savedBlogPosts, setSavedBlogPosts] = useState([]);
  const addAndSaveButtonOnPress = async (_: PressEvent) => {
    const payload: Array<BlogPost> = Object.assign([], scrapedBlogPosts);
    payload.push({
      author: newBlogPostAuthor.value,
      date: newBlogPostDate.value,
      title: newBlogPostTitle.value,
      url: newBlogPostURL.value,
    });

    const resp = await fetch('http://localhost:8000/api/persist-first', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(payload),
    });
    console.log("persist first response", resp);
    const data = await resp.json();
    setSavedBlogPosts(data);
  };

  const [savedBlogPosts2, setSavedBlogPosts2] = useState([]);
  const addAndSaveButtonOnPress2 = async (_: PressEvent) => {
    const payload: Array<BlogPost> = Object.assign([], savedBlogPosts);

    const resp = await fetch('http://localhost:8000/api/persist-second', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(payload),
    });

    console.log("persist second response", resp);

    const data = await resp.json();
    setSavedBlogPosts2(data);
  };

  const deleteBlogPostsByAuthorInput = useInput("");
  const [remainingBlogPosts, setRemainingBlogPosts] = useState([]);
  const deleteBlogPostsByAuthorButtonOnPress = async (_: PressEvent) => {
    const resp = await fetch('http://localhost:8000/api/delete?' + new URLSearchParams({
      author: deleteBlogPostsByAuthorInput.value,
    }).toString(), {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    console.log("delete response", resp);

    const data = await resp.json();
    setRemainingBlogPosts(data);
  };

  return (
    <main>
      <Button onPress={scrapeButtonOnPress}>Scrape Drew DeVault's blog</Button>

      {
        blogPostsTableFC('Scraped blog posts', scrapedBlogPosts)
      }

      <Spacer y={2} />

      <Input label="Author" placeholder="Author" {...newBlogPostAuthor.bindings}></Input>
      <Spacer y={1} />
      <Input label="Date" placeholder="Date" {...newBlogPostDate.bindings}></Input>
      <Spacer y={1} />
      <Input label="Title" placeholder="Title" {...newBlogPostTitle.bindings}></Input>
      <Spacer y={1} />
      <Input label="URL" placeholder="URL" {...newBlogPostURL.bindings}></Input>
      <Spacer y={1} />
      <Button onPress={addAndSaveButtonOnPress}>Add new blog post and save all</Button>

      {
        blogPostsTableFC('Saved blog posts', savedBlogPosts)
      }

      <Spacer y={2} />

      <Button onPress={addAndSaveButtonOnPress2}>Save all again</Button>

      {
        blogPostsTableFC('Saved blog posts 2', savedBlogPosts2)
      }

      <Spacer y={2} />

      <Input label="Author" placeholder="Author" {...deleteBlogPostsByAuthorInput.bindings}></Input>
      <Spacer y={1} />
      <Button onPress={deleteBlogPostsByAuthorButtonOnPress}>Delete blog posts by author</Button>

      {
        blogPostsTableFC('Remaining blog posts', remainingBlogPosts)
      }
    </main>
  )
}
