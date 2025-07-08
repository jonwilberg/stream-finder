import type { Route } from "./+types/home";
import { Welcome } from "../welcome/welcome";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Stream Finder" },
    { name: "description", content: "Find where to stream your favorite movies and TV shows" },
  ];
}

export default function Home() {
  return <Welcome />;
}
