#include "Player.hh"

/**
 * Write the name of your player and save this file
 * with the same name and .cc extension.
 */
#define PLAYER_NAME JohnnyMelavo
struct PLAYER_NAME : public Player
{
    /**
   * Factory: returns a new instance of this class.
   * Do not modify this function.
   */
    static Player *factory()
    {
        return new PLAYER_NAME;
    }

    /**
   * Types and attributes for your player can be defined here.
   */

    const vector<Dir> dirs = {Up, Down, Left, Right};
    // adalt, abaix, esquerra, dreta
    const vector<vector<int>> mov = {{-1, 0}, {1, 0}, {0, -1}, {0, 1}};

    map<Pos, Pos> Path;

    Pos BFS(Pos p, char objecte, bool fightmode)
    {
        queue<Pos> cola;
        cola.push(p);
        while (not cola.empty())
        {
            Pos pact = cola.front();
            cola.pop();
            if (pos_ok(pact))
            {
                Cell c = cell(pact);

                if (c.weapon == Bazooka and objecte == 'b')
                    return pact;
                else if (c.weapon == Gun and objecte == 'g')
                    return pact;
                else if (c.bonus == Money and objecte == 'm')
                    return pact;
                else if (c.id != -1 and objecte == 'B')
                {
                        return pact;
                }
               // else if (c.id != -1 and objecte == 'w')
               // {
               //     if (citizen(c.b_owner).type == Warrior and citizen(c.b_owner).weapon != Bazooka)
               //         return pact;
               // }

                for (int i = 0; i < 4; ++i)
                {
                    Pos ptemp = pact;
                    ptemp.i += mov[i][0];
                    ptemp.j += mov[i][1];

                    if (not fightmode)
                    {
                        if (pos_ok(ptemp) and cell(ptemp).type == Street and cell(ptemp).id == -1 and cell(ptemp).resistance == -1)
                        {
                            if (Path.find(ptemp) == Path.end())
                            {
                                Path[ptemp] = pact;
                                cola.push(ptemp);
                            }
                        }
                    }
                    else
                    {
                        if (pos_ok(ptemp) and cell(ptemp).type == Street and cell(ptemp).resistance == -1)
                        {
                            if (Path.find(ptemp) == Path.end())
                            {
                                Path[ptemp] = pact;
                                cola.push(ptemp);
                            }
                        }
                    }
                }
            }
        }
        return Pos(-1, -1);
    }
    Dir Pos2dir(Pos p)
    {
        if (p == Pos(-1, 0))
            return Up;
        else if (p == Pos(1, 0))
            return Down;
        else if (p == Pos(0, -1))
            return Left;
        else if (p == Pos(0, 1))
            return Right;
        else
        {
            assert(false); // this should never happen
            return Up;
        }
    }
    // Invariant: l'objecte ha d'existir en el mapa en el moment de la crida i l'hem trobat amb el BFS
    Dir direccionBFS(Pos pciti, Pos pobj)
    {
        Pos ptemp = pobj;
        vector<Pos> positions;
        while (pciti != ptemp)
        {
            positions.push_back(ptemp);
            ptemp = Path.at(ptemp);
        }
        Dir dr;
        if (not positions.empty())
        {
            Pos r;
            auto it = positions.end();
            --it;
            r.i = (*it).i - pciti.i;
            r.j = (*it).j - pciti.j;
            dr = Pos2dir(r);
        }
        else
        {
            assert(pciti == pobj);
            assert(positions.empty());
            dr = Up;
        }
        return dr;
    }

    /**
   * Play method, invoked once per each round.
   */
    virtual void play()
    {
        // INICIALITZAR INVARIANT
        Path.clear();
        // Suposem que existeixen bazokas pistolas i diners dins el mapa
        //bool Bexist = true;
        //bool wexist = true;
        //bool bexist = true;
        //bool gexist = true;
        //bool mexist = true;
        if (is_day())
        {
            vector<int> w = warriors(me());

            // Iterem sobre els guerrers
            for (int id : w)
            {
                Pos pciti = citizen(id).pos;
                Pos pobj;
                // Serà true si busquem amb Bfs (bazokas, pistoles...)
                bool moved = false;

                pobj = BFS(pciti, 'b', false); // busquem bazokas
                if (pobj != Pos(-1, -1) and pobj != pciti)
                {
                    Dir d = direccionBFS(pciti, pobj);
                    Path.clear();
                    if (pos_ok(pciti + d))
                    {
                        move(id, d);
                        moved = true;
                    }
                }
                else
                {
                    Path.clear();
                    //bexist = false; // no hi han bazokas al mapa
                }

                pobj = BFS(pciti, 'g', false); // busquem pistoles
                if (pobj != Pos(-1, -1) and pobj != pciti)
                {
                    Dir d = direccionBFS(pciti, pobj);
                    Path.clear();
                    if (pos_ok(pciti + d))
                    {
                        move(id, d);
                        moved = true;
                    }
                }
                else
                {
                    Path.clear();
                    //gexist = false; // no hi ha pistoles al mapa
                }

                pobj = BFS(pciti, 'm', false); // busquem diners
                if (pobj != Pos(-1, -1) and pobj != pciti)
                {
                    Dir d = direccionBFS(pciti, pobj);
                    Path.clear();
                    if (pos_ok(pciti + d))
                    {
                        move(id, d);
                        moved = true;
                    }
                }
                else
                {
                    Path.clear();
                    //mexist = false; // no hi ha diners al mapa
                }
            }
            Path.clear();
            vector<int> b = builders(me());
            for (int id : b)
            { // iterate over builders
                Pos pobj;
                Pos pciti = citizen(id).pos;
                bool moved = false;

                pobj = BFS(pciti, 'm', false); // busquem diners
                if (pobj != Pos(-1, -1) and pobj != pciti)
                {
                    Dir d = direccionBFS(pciti, pobj);
                    Path.clear();
                    if (pos_ok(pciti + d))
                    {
                        moved = true;
                        move(id, d);
                    }
                }
                else
                {
                    Path.clear();
                    //mexist = false; // no hi ha diners al mapa
                }
                pobj = BFS(pciti, 'b', false);
                if (pobj != Pos(-1, -1) and pobj != pciti)
                {
                    Dir d = direccionBFS(pciti, pobj);
                    Path.clear();
                    if (pos_ok(pciti + d))
                    {
                        moved = true;
                        move(id, d);
                    }
                }
                else
                {
                    Path.clear();
                    //bexist = false;
                }
                pobj = BFS(pciti, 'g', false);
                if (pobj != Pos(-1, -1) and pobj != pciti)
                {
                    Dir d = direccionBFS(pciti, pobj);
                    Path.clear();
                    if (pos_ok(pciti + d))
                    {
                        moved = true;
                        move(id, d);
                    }
                }
                else
                {
                    Path.clear();
                    //gexist = false;
                }
                Path.clear();
            }
        }
        else
        { // night time
            vector<int> w = warriors(me());
            for (int id : w)
            { // iterate over warriors
                Pos pciti = citizen(id).pos;
                Pos pobj;
                // Serà true si busquem amb Bfs (bazokas, pistoles...)
                bool moved = false;

                pobj = BFS(pciti, 'B', true); // busquem Builders
                if (pobj != Pos(-1, -1) and pobj != pciti)
                {
                    Dir d = direccionBFS(pciti, pobj);
                    Path.clear();
                    if (pos_ok(pciti + d))
                    {
                        move(id, d);
                        moved = true;
                    }
                }
                else
                {
                    Path.clear();
                }
            }

            vector<int> b = builders(me());
            for (int id : b)
            {
            }
        }
    }
};

/**
 * Do not modify the following line.
 */
RegisterPlayer(PLAYER_NAME);
