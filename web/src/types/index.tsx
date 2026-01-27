export type Repository = {
  id: number;
  name: string;
  owner: string;
  description?: string;
  language?: string;
  stars?: number;
  forks?: number;
  openIssues?: number;
  lastUpdated?: string; // ISO date string
  url?: string;
  private: boolean;
};

export type LanguageColor = {
  [language: string]: string; // Maps a language to its hex color
};

export type User = {
  id: string;
  github_id: string;
  username: string;
  email?: string;
  avatar_url?: string;
  access_token?: string;
  created_at: string;
  updated_at: string;
};
