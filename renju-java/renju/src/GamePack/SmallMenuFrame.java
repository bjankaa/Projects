package GamePack;

import MenuPack.MenuFrame;
import MenuPack.RuleFrame;

import javax.swing.*;

import java.awt.*;
import java.awt.event.ActionEvent;
import java.awt.event.ActionListener;

import static javax.swing.WindowConstants.EXIT_ON_CLOSE;

public class SmallMenuFrame extends JFrame {

    private GameFrame gamefr;

    public GameFrame getGameFrame(){ return gamefr; }
    public void setGameFrame(GameFrame gamefr ){ this.gamefr = gamefr; }

    private void initComponents(){
        this.setLayout(new GridLayout(3,1));

        JButton folytatas = new JButton("Folytatás");
        folytatas.setFont(new Font("Serif", Font.PLAIN, 20));
        this.add(folytatas);
        folytatas.addActionListener(new ActionListener() {
            @Override
            public void actionPerformed(ActionEvent e) {
                gamefr.setEnabled(true);
                setVisible(false);
            }
        });

        JButton newgame = new JButton ("Új játék");
        newgame.setFont(new Font("Serif", Font.PLAIN, 20));
        this.add(newgame);
        newgame.addActionListener(new ActionListener() {
            @Override
            public void actionPerformed(ActionEvent e) {
                setVisible(false);
                gamefr.close();
                GameFrame gf = new GameFrame();
                gf.run();
            }
        });


        JButton fooldal = new JButton("Vissza a főoldalra");
        fooldal.setFont(new Font("Serif", Font.PLAIN, 20));
        this.add(fooldal);
        fooldal.addActionListener(new ActionListener() {
            @Override
            public void actionPerformed(ActionEvent e) {
                setVisible(false);
                gamefr.close();
                MenuFrame mf = new MenuFrame();
                mf.run();
            }
        });

    }

    public SmallMenuFrame() {
        super("Pause");
        setDefaultCloseOperation(EXIT_ON_CLOSE);

        setSize(300, 200);
        setLocationRelativeTo(null);
        setResizable(false);
        initComponents();
    }
    public void run(){
        this.setVisible(true);
    }
}
