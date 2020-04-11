
// This file is part of the MDCII Game Engine.
// Copyright (C) 2020  Armin Schlegel
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.

#include <iostream>

#include "SDL2/SDL.h"

#include "sdlgui/layout.h"
#include "sdlgui/messagedialog.h"
#include "sdlgui/screen.h"
#include "sdlgui/texturebutton.h"
#include "sdlgui/textureview.h"
#include "sdlgui/window.h"

#include "bsh_texture.hpp"
#include "cod_parser.hpp"
#include "files.hpp"
#include "fps.hpp"
#include "gamewindow.hpp"
#include "haeuser.hpp"
#include "mainmenu.hpp"
#include "palette.hpp"

using namespace sdlgui;
MainMenu::MainMenu(
    SDL_Renderer* renderer, std::shared_ptr<Basegad> basegad, SDL_Window* pwindow, int rwidth, int rheight, bool fullscreen, const std::string& gam_name)
  : renderer(renderer)
  , basegad(basegad)
  , width(rwidth)
  , height(rheight)
  , fullscreen(fullscreen)
  , gam_name(gam_name)
  , pwindow(pwindow)
  , files(Files::instance())
  , triggerStartGame(false)
  , Screen(pwindow, Vector2i(rwidth, rheight), "Game", false, true)
{
  std::cout << "Basegad: " << basegad->get_gadgets_size() << std::endl;
  Bsh_leser bsh_leser(files->instance()->find_path_for_file("toolgfx/start.bsh"));
  BshImageToSDLTextureConverter converter(renderer);

  SDL_Texture* background = converter.Convert(&bsh_leser.gib_bsh_bild(0));

  auto singlePlayerButtonGad = basegad->get_gadgets_by_index(2);
  SDL_Texture* singlePlayerTexture = converter.Convert(&bsh_leser.gib_bsh_bild(singlePlayerButtonGad->Gfxnr));
  SDL_Texture* singlePlayerTextureClicked = converter.Convert(&bsh_leser.gib_bsh_bild(singlePlayerButtonGad->Gfxnr + singlePlayerButtonGad->Pressoff));

  auto multiPlayerButtonGad = basegad->get_gadgets_by_index(3);
  SDL_Texture* multiPlayerTexture = converter.Convert(&bsh_leser.gib_bsh_bild(multiPlayerButtonGad->Gfxnr));
  SDL_Texture* multiPlayerTextureClicked = converter.Convert(&bsh_leser.gib_bsh_bild(multiPlayerButtonGad->Gfxnr + multiPlayerButtonGad->Pressoff));

  auto creditsButtonGad = basegad->get_gadgets_by_index(4);
  SDL_Texture* creditsTexture = converter.Convert(&bsh_leser.gib_bsh_bild(creditsButtonGad->Gfxnr));
  SDL_Texture* creditsTextureClicked = converter.Convert(&bsh_leser.gib_bsh_bild(creditsButtonGad->Gfxnr + creditsButtonGad->Pressoff));

  auto introButtonGad = basegad->get_gadgets_by_index(5);
  SDL_Texture* introTexture = converter.Convert(&bsh_leser.gib_bsh_bild(introButtonGad->Gfxnr));
  SDL_Texture* introTextureClicked = converter.Convert(&bsh_leser.gib_bsh_bild(introButtonGad->Gfxnr + introButtonGad->Pressoff));

  auto exitButtonGad = basegad->get_gadgets_by_index(6);
  SDL_Texture* exitTexture = converter.Convert(&bsh_leser.gib_bsh_bild(exitButtonGad->Gfxnr));
  SDL_Texture* exitTextureClicked = converter.Convert(&bsh_leser.gib_bsh_bild(exitButtonGad->Gfxnr + exitButtonGad->Pressoff));

  {
    wdg<TextureView>(background);
    auto& singlePlayerButton = wdg<TextureButton>(singlePlayerTexture, [this] {
      std::cout << "Singleplayer pressed" << std::endl;
      triggerStartGame = true;
    });
    singlePlayerButton.setPosition(singlePlayerButtonGad->Pos.x, singlePlayerButtonGad->Pos.y);
    singlePlayerButton.setSecondaryTexture(singlePlayerTextureClicked);
    singlePlayerButton.setTextureSwitchFlags(TextureButton::OnClick);

    auto& multiPlayerButton = wdg<TextureButton>(multiPlayerTexture, [this] { std::cout << "Multiplayer pressed" << std::endl; });
    multiPlayerButton.setPosition(multiPlayerButtonGad->Pos.x, multiPlayerButtonGad->Pos.y);
    multiPlayerButton.setSecondaryTexture(multiPlayerTextureClicked);
    multiPlayerButton.setTextureSwitchFlags(TextureButton::OnClick);

    auto& creditsButton = wdg<TextureButton>(creditsTexture, [this] { std::cout << "credits pressed" << std::endl; });
    creditsButton.setPosition(creditsButtonGad->Pos.x, creditsButtonGad->Pos.y);
    creditsButton.setSecondaryTexture(creditsTextureClicked);
    creditsButton.setTextureSwitchFlags(TextureButton::OnClick);

    auto& introButton = wdg<TextureButton>(introTexture, [this] { std::cout << "intro pressed" << std::endl; });
    introButton.setPosition(introButtonGad->Pos.x, introButtonGad->Pos.y);
    introButton.setSecondaryTexture(introTextureClicked);
    introButton.setTextureSwitchFlags(TextureButton::OnClick);

    auto& exitButton = wdg<TextureButton>(exitTexture, [this] {
      std::cout << "exit pressed" << std::endl;
      exit(0);
    });
    exitButton.setPosition(exitButtonGad->Pos.x, exitButtonGad->Pos.y);
    exitButton.setSecondaryTexture(exitTextureClicked);
    exitButton.setTextureSwitchFlags(TextureButton::OnClick);
  }
  performLayout(mSDL_Renderer);
}

void MainMenu::LoadGame(const std::string& gam_name)
{
  auto haeuser_cod = std::make_shared<Cod_Parser>(files->instance()->find_path_for_file("haeuser.cod"), true, false);
  auto haeuser = std::make_shared<Haeuser>(haeuser_cod);

  GameWindow gameWindow(renderer, haeuser, pwindow, width, height, gam_name, fullscreen);
  gameWindow.Handle();
  Handle();
}

// todo: add signal/slot for exiting window
void MainMenu::Handle()
{
  auto palette = Palette::instance();

  auto transparentColor = palette->getColor(palette->getTransparentColor());
  std::cout << "[INFO] Transparent color: " << (int)transparentColor.getRed() << ", " << (int)transparentColor.getGreen() << ", "
            << (int)transparentColor.getBlue() << std::endl;

  SDL_Texture* texture;
  SDL_Surface* final_surface;

  SDL_Surface* s8 = SDL_CreateRGBSurface(0, width, height, 8, 0, 0, 0, 0);
  SDL_Color c[palette->size()];
  int i, j;
  for (i = 0, j = 0; i < palette->size(); i++)
  {
    c[i].r = palette->getColor(i).getRed();
    c[i].g = palette->getColor(i).getGreen();
    c[i].b = palette->getColor(i).getBlue();
  }
  SDL_SetPaletteColors(s8->format->palette, c, 0, palette->size());

  try
  {
    final_surface = SDL_ConvertSurfaceFormat(s8, SDL_PIXELFORMAT_RGB888, 0);
    texture = SDL_CreateTextureFromSurface(renderer, final_surface);
    SDL_RenderClear(renderer);
    SDL_RenderCopy(renderer, texture, NULL, NULL);
    this->drawAll();
    SDL_RenderPresent(renderer);

    Fps fps;
    const Uint8* keystate = SDL_GetKeyboardState(NULL);

    SDL_Event e;
    bool quit = false;
    while (!quit)
    {
      while (SDL_PollEvent(&e) != 0)
      {
        switch (e.type)
        {
          case SDL_QUIT:
            quit = true;
            break;
          case SDL_USEREVENT:
            break;
          case SDL_KEYDOWN:
            if (e.key.keysym.sym == SDLK_ESCAPE)
            {
              quit = true;
            }
            break;
        }
        this->onEvent(e);
      }
      int x, y;
      SDL_GetMouseState(&x, &y);

      final_surface = SDL_ConvertSurfaceFormat(s8, SDL_PIXELFORMAT_RGB888, 0);
      texture = SDL_CreateTextureFromSurface(renderer, final_surface);
      SDL_FreeSurface(final_surface);
      SDL_RenderClear(renderer);
      SDL_RenderCopy(renderer, texture, NULL, NULL);
      SDL_SetTextureBlendMode(texture, SDL_BLENDMODE_NONE);
      if (triggerStartGame)
      {
        triggerStartGame = false;
        this->LoadGame(this->gam_name);
      }
      this->drawAll();
      SDL_RenderPresent(renderer);
      SDL_DestroyTexture(texture);
      fps.next();
    }
  }
  catch (const std::runtime_error& e)
  {
    std::string error_msg = std::string("Caught a fatal error: ") + std::string(e.what());
#if defined(_WIN32)
    MessageBoxA(nullptr, error_msg.c_str(), NULL, MB_ICONERROR | MB_OK);
#else
    std::cerr << error_msg << std::endl;
#endif
  }
  SDL_FreeSurface(s8);
}
