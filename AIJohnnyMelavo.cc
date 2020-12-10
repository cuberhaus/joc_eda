#include "Player.hh"

/**
 * Write the name of your player and save this file
 * with the same name and .cc extension.
 */
#define PLAYER_NAME JohnnyBGood
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
        int player = citizen(cell(p).id).player;
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
                    if (citizen(c.id).type == Builder and citizen(c.id).player != player)
                    {
                        return pact;
                    }
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
        if (is_day())
        {
            vector<int> w = warriors(me());
            // Iterem sobre els guerrers
            for (int id : w)
            {
                Pos pciti = citizen(id).pos;
                Pos pobj;
                if (citizen(cell(pciti).id).weapon != Bazooka)
                {
                    // busquem bazookas
                    Path.clear();
                    pobj = BFS(pciti, 'b', false);
                    if (pobj != Pos(-1, -1) and pobj != pciti)
                    {
                        Dir d = direccionBFS(pciti, pobj);
                        if (pos_ok(pciti + d))
                        {
                            move(id, d);
                        }
                    }
                    // busquem pistoles
                    Path.clear();
                    pobj = BFS(pciti, 'g', false);
                    if (pobj != Pos(-1, -1) and pobj != pciti)
                    {
                        Dir d = direccionBFS(pciti, pobj);
                        if (pos_ok(pciti + d))
                        {
                            move(id, d);
                        }
                    }
                }

                // busquem diners
                Path.clear();
                pobj = BFS(pciti, 'B', false);
                if (pobj != Pos(-1, -1) and pobj != pciti)
                {
                    Dir d = direccionBFS(pciti, pobj);
                    if (pos_ok(pciti + d))
                    {
                        move(id, d);
                    }
                }
            }
            vector<int> b = builders(me());
            for (int id : b)
            { // iterate over builders
                Pos pobj;
                Pos pciti = citizen(id).pos;
                // busquem diners
                Path.clear();
                pobj = BFS(pciti, 'm', false);
                if (pobj != Pos(-1, -1) and pobj != pciti)
                {
                    Dir d = direccionBFS(pciti, pobj);
                    if (pos_ok(pciti + d))
                    {
                        move(id, d);
                    }
                }
                // Busquem bazookas
                Path.clear();
                pobj = BFS(pciti, 'b', false);
                if (pobj != Pos(-1, -1) and pobj != pciti)
                {
                    Dir d = direccionBFS(pciti, pobj);
                    if (pos_ok(pciti + d))
                    {
                        move(id, d);
                    }
                }
                // Busquem pistoles
                Path.clear();
                pobj = BFS(pciti, 'g', false);
                if (pobj != Pos(-1, -1) and pobj != pciti)
                {
                    Dir d = direccionBFS(pciti, pobj);
                    Path.clear();
                    if (pos_ok(pciti + d))
                    {
                        move(id, d);
                    }
                }
            }
        }
        else
        { // night time
            vector<int> w = warriors(me());
            for (int id : w)
            { // iterate over warriors
                Pos pciti = citizen(id).pos;
                Pos pobj;

                // busquem Builders
                Path.clear();
                pobj = BFS(pciti, 'B', true);
                if (pobj != Pos(-1, -1) and pobj != pciti)
                {
                    Dir d = direccionBFS(pciti, pobj);
                    if (pos_ok(pciti + d))
                    {
                        move(id, d);
                    }
                }
            }

            vector<int> b = builders(me());
            for (int id : b)
            {
                Pos pobj;
                Pos pciti = citizen(id).pos;
                // busquem diners
                Path.clear();
                pobj = BFS(pciti, 'm', false);
                if (pobj != Pos(-1, -1) and pobj != pciti)
                {
                    Dir d = direccionBFS(pciti, pobj);
                    if (pos_ok(pciti + d))
                    {
                        move(id, d);
                    }
                }
                // busquem bazookas
                Path.clear();
                pobj = BFS(pciti, 'b', false);
                if (pobj != Pos(-1, -1) and pobj != pciti)
                {
                    Dir d = direccionBFS(pciti, pobj);
                    if (pos_ok(pciti + d))
                    {
                        move(id, d);
                    }
                }
                // busquem pistoles
                Path.clear();
                pobj = BFS(pciti, 'g', false);
                if (pobj != Pos(-1, -1) and pobj != pciti)
                {
                    Dir d = direccionBFS(pciti, pobj);
                    if (pos_ok(pciti + d))
                    {
                        move(id, d);
                    }
                }
            }
        }
    }
};

/**
 * Do not modify the following line.
 */
RegisterPlayer(PLAYER_NAME);
