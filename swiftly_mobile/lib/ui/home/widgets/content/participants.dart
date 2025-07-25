import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:swiftly_mobile/providers/user_notifier_provider.dart';
import 'package:swiftly_mobile/ui/core/themes/theme.dart';
import 'package:swiftly_mobile/ui/home/widgets/content/user_widget.dart';

class Participants extends ConsumerWidget {
  const Participants({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final users = ref.watch(userNotifierProvider).users.toList();
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text('Участники команды', style: AppTextStyles.style12),
        const SizedBox(height: 10),
        SingleChildScrollView(
          scrollDirection: Axis.horizontal,
          child: Row(
            children: [
              ...users
                  .expand(
                    (user) => [
                      const SizedBox(width: 5),
                      UserWidget(user: user),
                    ],
                  )
                  .skip(1),
            ],
          ),
        ),
      ],
    );
  }
}
