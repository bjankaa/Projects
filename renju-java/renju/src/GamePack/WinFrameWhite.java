package GamePack;

import MenuPack.MenuFrame;

import javax.swing.*;
import java.awt.*;
import java.awt.event.ActionEvent;
import java.awt.event.ActionListener;
import java.util.Objects;

public class WinFrameWhite extends JFrame {

    private GameFrame gamefr;

    public void setGameFrame(GameFrame gamefr ){ this.gamefr = gamefr; }

    public void initComponents(){
        this.setLayout(new GridLayout(3,1));


        JLabel kiiras = new JLabel("A győztes fehér");
        kiiras.setHorizontalAlignment(JLabel.CENTER);//what if dontetlen
        kiiras.setFont(new Font("Serif", Font.PLAIN, 20));
        this.add(kiiras);


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

    public WinFrameWhite() {
        super("Fehér Győzelem");
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

