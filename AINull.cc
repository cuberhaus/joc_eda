#include "Player.hh"

/**
 * Write the name of your player and save this file
 * with the same name and .cc extension.
 */
#define PLAYER_NAME Null
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

    Pos BFS(Pos p, char objecte)
    {
        //vector<vector<bool>> visitats(board_cols(), vector<bool>(board_rows(), false));
        queue<Pos> cola;
        cola.push(p);
        while (not cola.empty())
        {
            // p has to be the node visited
            Pos pact = cola.front();
            cola.pop();
            if (pos_ok(pact))
            {
                Cell c = cell(pact);
                //visitats[pact.i][pact.j] = true;

                if (c.weapon == Bazooka and objecte == 'b')
                    return pact;
                else if (c.weapon == Gun and objecte == 'g')
                    return pact;

                for (int i = 0; i < 4; ++i)
                {
                    Pos ptemp = pact;
                    ptemp.i += mov[i][0];
                    ptemp.j += mov[i][1];
                    if (pos_ok(ptemp) and cell(ptemp).type == Street and cell(ptemp).id == -1) //and not visitats[ptemp.i][ptemp.j])
                    {
                        if (Path.find(ptemp) == Path.end()) {
                            Path[ptemp] = pact;
                            cola.push(ptemp);
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
        Pos r;
        auto it = positions.end();
        --it;
        r.i = (*it).i - pciti.i;
        r.j = (*it).j - pciti.j;
        //cerr << r;
        Dir dr = Pos2dir(r);
        return dr;
    }

    /**
   * Play method, invoked once per each round.
   */
    virtual void play()
    {
        // if (start) {
        //     start = false;
        //     bazookas.resize(num_ini_bazookas());
        //     guns.resize(num_ini_guns());
        //     moneys.resize(num_ini_money());
        //     foods.resize(num_ini_food());
        //     for (int i = 0; i < num_ini_bazookas(); ++i) {
        //         bazookas[i].i = -1;
        //         bazookas[i].j = -1;
        //     }
        //     for (int i = 0; i < num_ini_guns(); ++i) {
        //         guns[i].i = -1;
        //         guns[i].j = -1;
        //     }
        //     for (int i = 0; i < num_ini_money(); ++i) {
        //         moneys[i].i = -1;
        //         moneys[i].j = -1;
        //     }
        //     for (int i = 0; i < num_ini_food(); ++i) {
        //         foods[i].i = -1;
        //         foods[i].j = -1;
        //     }
        // }

        // n = board max column, m = board max row
        // int n = board_cols();
        // int m = board_rows();
        // for (int i = 0; i < n; ++i)
        // {
        //     for (int j = 0; j < m; ++j)
        //     {
        //         Cell casilla = cell(i, j);
        //         if (casilla.weapon == Bazooka) {
        //
        //         }
        //         cerr << casilla.type << ' '; // cout << street or building
        //     }
        //     cerr << endl;
        // }

        // INICIALITZAR INVARIANTS
        Path.clear();
        //if (is_day())
        //{
        vector<int> w = warriors(me());
        bool bexist = true;
        bool gexist = true;
        bool mexist = true;
        bool bdone = false;
        bool gdone = false;
        bool mdone = false;
        for (int id : w)
        { // iterate over warriors
            Pos pciti = citizen(id).pos;
            bdone = false;
            gdone = false;
            mdone = false;
            if (bexist)
            {
                Pos pobj = BFS(pciti, 'b'); // buscamos bazokas
                bdone = true;
                if (pobj != Pos(-1, -1))
                {
                    Dir d = direccionBFS(pciti, pobj);
                    if (pos_ok(pciti + d))
                    {
                        move(id, d);
                    }
                }
                else
                {
                    bexist = false;
                }
            }
            else if (gexist)
            {
                Pos pobj = BFS(pciti, 'g');
                gdone = true;
                if (pobj != Pos(-1, -1))
                {
                    Dir d = direccionBFS(pciti, pobj);
                    if (pos_ok(pciti + d))
                    {
                        move(id, d);
                    }
                }
                else
                {
                    gexist = false;
                }
            }
            Path.clear();
        }
            //}

            //     // At day take care of builders
            vector<int> b = builders(me());
            for (int id : b)
            { // iterate over builders

                Pos pciti = citizen(id).pos;
                bdone = false;
                gdone = false;
                mdone = false;
                if (bexist)
                {
                    Pos pobj = BFS(pciti, 'b');
                    bdone = true;
                    if (pobj != Pos(-1, -1))
                    {
                        Dir d = direccionBFS(pciti, pobj);
                        if (pos_ok(pciti + d))
                        {
                            move(id, d);
                        }
                    }
                    else
                    {
                        bexist = false;
                    }
                }
                else if (gexist)
                {
                    Pos pobj = BFS(pciti, 'g');
                    gdone = true;
                    if (pobj != Pos(-1, -1))
                    {
                        Dir d = direccionBFS(pciti, pobj);
                        if (pos_ok(pciti + d))
                        {
                            move(id, d);
                        }
                    }
                    else
                    {
                        gexist = false;
                    }
                }
            Path.clear();
            }
            // }
            // else
            // { // night time
            //     vector<int> w = warriors(me());
            //     for (int id : w)
            //     { // iterate over warriors
            //         Pos pciti = citizen(id).pos;
            //         //cerr << "citizen " << id << " in position: " << pciti << endl;
            //         char c = 'b'; // busquem bazokas
            //         Pos pobj = BFS(pciti, c);
            //         if (pobj != Pos(-1, -1))
            //         {
            //             Dir d = direccionBFS(pciti, pobj);
            //             //Dir d = Down; // prova
            //             if (pos_ok(pciti + d))
            //             {
            //                 move(id, d);
            //             }
            //         }
            //     }

            //     // At day take care of builders
            //     vector<int> b = builders(me());
            //     for (int id : b)
            //     { // iterate over builders
            //         Pos pciti = citizen(id).pos;
            //         //cerr << "citizen " << id << " in position: " << pciti << endl;
            //         char c = 'b'; // busquem bazokas
            //         Pos pobj = BFS(pciti, c);
            //         if (pobj != Pos(-1, -1))
            //         {
            //             Dir d = direccionBFS(pciti, pobj);
            //             //Dir d = Down; // prova
            //             if (pos_ok(pciti + d))
            //             {
            //                 move(id, d);
            //             }
            //         }
            //     }
        
    }
};

/**
 * Do not modify the following line.
 */
RegisterPlayer(PLAYER_NAME);
